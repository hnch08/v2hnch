export namespace config {
	
	export class Config {
	    username: string;
	    name: string;
	    requestURL: string;
	    port: string;
	    address: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.name = source["name"];
	        this.requestURL = source["requestURL"];
	        this.port = source["port"];
	        this.address = source["address"];
	        this.id = source["id"];
	    }
	}

}

