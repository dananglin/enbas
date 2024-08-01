# Configuration reference

## Config

| Field              | Type                          | Description                                      |
|--------------------|-------------------------------|--------------------------------------------------|
| `credentialsFile`  | string                        | The (absolute) path to your credentials file.    |
| `cacheDirectory`   | string                        | The (absolute) path to the root cache directory. |
| `lineWrapMaxWidth` | int                           | The character limit used for line wrapping.      |
| `http`             | [HTTPConfig](#httpconfig)     | HTTP settings.                                   |
| `integrations`     | [Integrations](#integrations) | Specify your integrations with Enbas.            |


## HTTPConfig

| Field          | Type | Description                                                         |
|----------------|------|---------------------------------------------------------------------|
| `timeout`      | int  | The timeout (in seconds) for normal HTTP requests to your instance. |
| `mediaTimeout` | int  | The timeout (in seconds) for retrieving media from your instance.   |

## Integrations

| Field         | Type   | Description                                                                                            |
|---------------|--------|--------------------------------------------------------------------------------------------------------|
| `browser`     | string | The browser used for opening URLs (e.g. URL to a remote account).                                      |
| `editor`      | string | The text editor used for writing statuses (not yet implemented).                                       |
| `pager`       | string | The pager used for paging through long outputs (e.g. status timelines). Leave blank to disable paging. |
| `imageViewer` | string | The image viewer used for viewing images from media attachments.                                       |
| `videoPlayer` | string | The video player used for viewing videos from media attachments.                                       |
