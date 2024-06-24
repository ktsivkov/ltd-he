export namespace player {
	
	export class Player {
	    battleTag: string;
	    logsDirPath: string;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battleTag = source["battleTag"];
	        this.logsDirPath = source["logsDirPath"];
	    }
	}

}

