import yaml from 'js-yaml';

class FabricClient {
    constructor(config) {
        this.connectionProfile = yaml.safeLoad(config);
    }

    login(username, password) {
        
    }
}

export default FabricClient;