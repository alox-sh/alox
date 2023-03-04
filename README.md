
# Alox

## Progress

### 

|Interface|Implementation|Status|
|-|-|-|
|Server|Server *(base)*|:heavy_check_mark:|
|^|API|:heavy_check_mark:|
|^|Website|:heavy_check_mark:|
|^|File|:white_check_mark:|
|^|Proxy|:white_check_mark:|
|Response|*(base)*|:heavy_check_mark:|
|^|JSON|:heavy_check_mark:|
|^|Page|:heavy_check_mark:|

<!-- - [-] Server
    - [x] API server
    - [x] Web server
    - [ ] File server
    - [ ] Proxy server -->

# Alox (client)

## Usage

1. Install with<br/>
    **NPM**<br/>
    `npm i --save @alox.sh/client`<br/>
    or **YARN**<br/>
    `yarn add @alox.sh/client`<br/>
2. Import
    ```ts
    import type { WebsiteState } from '@alox.sh/client'

    // Declare 
    export declare global {
        interface Window {
            __websiteState?: WebsiteState
        }
    }
    ```

## Supported application types

### Out of the box support
- SPA (served from file system) - see `examples/spaDir`
- SPA (embeded into the binary) - see `examples/spaEmbeded`

### Extendable
- Any other type - this library is fully customizable to support virtually any application type
