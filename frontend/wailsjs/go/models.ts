export namespace plugins {
	
	export class Info {
	    name: string;
	    description: string;
	    author: string;
	    version: string;
	    pluginId: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.version = source["version"];
	        this.pluginId = source["pluginId"];
	    }
	}

}

