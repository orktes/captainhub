function lc(str) {
  return str && str.toLowerCase ? str.toLowerCase() : str;
}

function pullRequestOpened(eventData) {
  var pullRequestUser = lc(eventData.pull_request.user.login);
  var reviewerCount = config.reviewerCount;
  var prNumber = eventData.number;

  var message = 'Awesome work! Now just sit back and wait for Travis to pass and `' + reviewerCount + '`` others to review your code.\n\n';

  message += '\n### Review commands\n';
  message += '- accept: `pr_review OK`\n';
  message += '- print status: `pr_review status`\n';


  createIssueComment(prNumber, message);
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
    if (message.indexOf('pr_review') === 0) {
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
          var reviewerCountStr = loadData(prNumber + ':waiting_review_count') || config.reviewerCount.toString();
          var reviewedByStr = loadData(prNumber + ':waiting_reviewed_by') || '[]';

          var reviewerCount = JSON.parse(reviewerCountStr);
          var reviewedBy = JSON.parse(reviewedByStr);

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
            createIssueComment(prNumber, 'Review completed!');
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
      }
    }
  });
}


// Process hook data
if (eventType === 'pull_request' && eventData.action === 'opened') {
  pullRequestOpened(eventData);
} else if (eventType === 'issue_comment' && eventData.action === 'created' && eventData.issue.pull_request) {
  pullRequestComment(eventData);
}
