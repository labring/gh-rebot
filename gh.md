
https://cli.github.com/manual/gh_issue_edit

gh pr checks 375 --repo cuisongliu/sealos
gh pr close  375 --repo cuisongliu/sealos
gh pr ready  375 --undo --repo cuisongliu/sealos
gh pr review  375 --repo cuisongliu/sealos
gh pr review  375 --repo cuisongliu/sealos --comment/--request-changes/--approve -b "interesting"

gh label create bug --description "Something isn't working" --color E99695  --repo cuisongliu/sealos
gh pr edit  375 --add-label "bug" --add-reviewer "cuisongliu" --add-assignee "cuisongliu"  --repo cuisongliu/sealos
gh run rerun 3310817036 --failed --repo cuisongliu/sealos

${{ github.event.comment.user.login }}
${{ github.event.sender.login }}
