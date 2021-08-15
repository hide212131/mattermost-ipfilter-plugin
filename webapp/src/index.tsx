import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import manifest from './manifest';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

const getPluginServerRoute = (state: GlobalState) => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + manifest.id;
};

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(
        registry: PluginRegistry,
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        store: Store<GlobalState, Action<Record<string, unknown>>>,
    ) {
        const route = getPluginServerRoute(store.getState());
        const response = await fetch(route);
        const text = await response.text();

        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        registry.registerFilesWillUploadHook((files, upload) => {
            if (text === 'OK') {
                return (files.length >= 2) ? {message: 'Must upload one by one.'} : {files};
            }
            return {message: 'Your host is not permitted'};
        });
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
