# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Remote Procedure Calls Golang Changelog

## <a name="1.0.7"></a> 1.0.7 (2023-03-29)

### Bug fixes
- Added error return when call closed client
## <a name="1.0.6"></a> 1.0.6 (2022-12-08)
### Bug fixes

- **clients** fixed HandleHttpResponse for nil values
## <a name="1.0.5"></a> 1.0.5 (2022-11-18)

### Bug fixes
- **services** fixed call of **CallCommand** with nil parameters

## <a name="1.0.4"></a> 1.0.4 (2022-10-27)

### Features

- **auth** replace context auth info strings on **AuthField**

## <a name="1.0.3"></a> 1.0.3 (2022-10-07)

### Bug fixes

- Fixed embedding
- Updated about, heartbeat and rest operations
- Updated CommandableSwaggerDocument objects conversion
- Updated dependencies

## <a name="1.0.2"></a> 1.0.2 (2022-07-05)

- Update dependencies

## <a name="1.0.0"></a> 1.0.0 (2022-06-24) 

Initial public release

### Features

* **build** HTTP service factory
* **clients** mechanisms for retrieving connection settings
* **connect** helper module to retrieve connections services and clients
* **services** basic implementation of services for connecting

