export namespace main {
	
	export class ConnectionConfig {
	    protocol: string;
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	    exchange: string;
	    queue: string;
	    routingKey: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.exchange = source["exchange"];
	        this.queue = source["queue"];
	        this.routingKey = source["routingKey"];
	        this.role = source["role"];
	    }
	}
	export class MessageItem {
	    protocol: string;
	    content: string;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new MessageItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.content = source["content"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PublishConfig {
	    protocol: string;
	    topic: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new PublishConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.topic = source["topic"];
	        this.message = source["message"];
	    }
	}
	export class SubscribeConfig {
	    protocol: string;
	    topic: string;
	
	    static createFrom(source: any = {}) {
	        return new SubscribeConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.topic = source["topic"];
	    }
	}

}

