# Attachment Filter Plugin

This plugin provides a filter function to allow or deny about attempting to upload attachments.

## Usage
1. Your server must have Enable set to true in the `File Sharing and Downloads > Allow File Sharing` section of its System Console.
2. After insatall plugin, Go to `Plugin > Attachment Filter Plugin` section. and describe it in JSON logic so that it returns true for the conditions of the files you want to allow. See [attachment_policy_test.json](https://github.com/hide212131/mattermost-plugin-attachment-filter/blob/master/server/attachment_policy_test.json) for information on what can be written in JSON logic.