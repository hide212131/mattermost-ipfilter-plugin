import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import manifest from './manifest';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(
        registry: PluginRegistry,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        store: Store<GlobalState, Action<Record<string, unknown>>>,
    ) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        registry.registerFilesWillUploadHook((files, upload) => {
            return (files.length >= 2) ? {message: 'Must upload one by one.'} : {files};
        });
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
