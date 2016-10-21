function lc(str) {
  return str && str.toLowerCase ? str.toLowerCase() : str;
}

function pullRequestOpened(eventData, update) {
  var pullRequestUser = lc(eventData.pull_request.user.login);
  var reviewerCount = config.reviewerCount;
  var prNumber = eventData.number;
  var message;

  if (!update) {
    message = 'Awesome work! Now just sit back and wait for Travis to pass and `' + reviewerCount + '` others to review your code.\n\n';

    message += '\n### Review commands\n';
    message += '- accept: `pr_review OK`\n';
    message += '- print status: `pr_review status`\n';
  } else {
    message = 'Some of the files were updated! `' + reviewerCount + '` others need to review the changes.';
  }

  if (!config.silent) {
    createIssueComment(prNumber, message);
  }
  createStatus(
    eventData.pull_request.head.sha,
    'pending',
    eventData.pull_request.url,
    'Pull request review waiting ' + reviewerCount + ' peers',
    'pr_review'
  );
  saveData(prNumber + ':waiting_review_count', JSON.stringify(reviewerCount));
  saveData(prNumber + ':waiting_reviewed_by', JSON.stringify([]));
}


function pullRequestComment(eventData) {
  var prNumber = eventData.issue.number;
  var senderName = lc(eventData.sender.login);
  var body = eventData.comment.body.trim();
  var cmd;

  _.each(body.split('\n'), function (message) {
    if (message.substr(0, 'pr_review'.length).toLowerCase() === 'pr_review') {
      cmd = _.filter(
        _.map(message.substring(10).split(' '), function (part) {
          return part.trim().replace(/^@/, '');
        }),
        function (cmd) {
          return !!cmd;
        }
      );

      console.log('cmd ' + cmd.join(' '));

      switch(cmd[0].toLowerCase()) {
        case 'ok':

          var prDetails = getPullRequestDetails(prNumber);
          var pullRequestUser = lc(prDetails.user.login);
          var reviewerCountStr = loadData(prNumber + ':waiting_review_count') || config.reviewerCount.toString();
          var reviewedByStr = loadData(prNumber + ':waiting_reviewed_by') || '[]';

          var reviewerCount = JSON.parse(reviewerCountStr);
          var reviewedBy = JSON.parse(reviewedByStr);

          if (senderName === pullRequestUser) {
            createIssueComment(prNumber, 'Cant review own pull request');
            return;
          }

          if (reviewedBy.indexOf(senderName) > -1) {
            // already reviewed by this user
            createIssueComment(prNumber, '@' + senderName + ' already reviewed!');
            return;
          }

          reviewedBy.push(senderName);
          reviewerCount--;

          if (reviewerCount > 0) {
            createStatus(
              prDetails.head.sha,
              'pending',
              prDetails.url,
              'Pull request review waiting ' + reviewerCount + ' peers',
              'pr_review'
            );
          } else {
            createStatus(
              prDetails.head.sha,
              'success',
              prDetails.url,
              'Pull request review done',
              'pr_review'
            );
            if (!config.silent) {
              createIssueComment(prNumber, 'Review completed!');
            }
          }

          saveData(prNumber + ':waiting_review_count', JSON.stringify(reviewerCount));
          saveData(prNumber + ':waiting_reviewed_by', JSON.stringify(reviewedBy));

        break;
        case 'status':
          var reviewerCountStr = loadData(prNumber + ':waiting_review_count') || config.reviewerCount.toString();
          var reviewedByStr = loadData(prNumber + ':waiting_reviewed_by') || '[]';

          var reviewerCount = JSON.parse(reviewerCountStr);
          var reviewedBy = JSON.parse(reviewedByStr);

          var message = 'Waiting review from `' + reviewerCount + ' peers`.\n\n';
          message += '\n### Reviewed by\n';

          _.each(reviewedBy, function (reviewer) {
              message += '- ' + reviewer;
          });

          createIssueComment(prNumber, message);
        break;
        case 'reopen':
          var fakeEvent = {number: prNumber, pull_request: getPullRequestDetails(prNumber)};
          pullRequestOpened(fakeEvent);
        break;
      }
    }
  });
}

// Process hook data
if (eventType === 'pull_request' && eventData.action === 'opened') {
  pullRequestOpened(eventData);
} else if (eventType === 'issue_comment' && eventData.action === 'created' && eventData.issue.pull_request) {
  pullRequestComment(eventData);
} else if (eventType === 'pull_request' && eventData.action === 'synchronize') {
  pullRequestOpened(eventData, true);
}
