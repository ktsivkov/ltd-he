export namespace history {
	
	export class GameHistory {
	    outcome: string;
	    eloDiff: number;
	    date: string;
	    gameId: number;
	    file: string;
	    totalGames: number;
	    wins: number;
	    elo: number;
	    totalLosses: number;
	    gamesLeftEarly: number;
	    winsStreak: number;
	    highestWinStreak: number;
	    mvp: number;
	    token: string;
	    player: string;
	    // Go type: time
	    timestamp: any;
	    payload: number[];
	
	    static createFrom(source: any = {}) {
	        return new GameHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outcome = source["outcome"];
	        this.eloDiff = source["eloDiff"];
	        this.date = source["date"];
	        this.gameId = source["gameId"];
	        this.file = source["file"];
	        this.totalGames = source["totalGames"];
	        this.wins = source["wins"];
	        this.elo = source["elo"];
	        this.totalLosses = source["totalLosses"];
	        this.gamesLeftEarly = source["gamesLeftEarly"];
	        this.winsStreak = source["winsStreak"];
	        this.highestWinStreak = source["highestWinStreak"];
	        this.mvp = source["mvp"];
	        this.token = source["token"];
	        this.player = source["player"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.payload = source["payload"];
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

}

export namespace player {
	
	export class Player {
	    battleTag: string;
	    logsDirPath: string;
	    reportFilePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battleTag = source["battleTag"];
	        this.logsDirPath = source["logsDirPath"];
	        this.reportFilePath = source["reportFilePath"];
	    }
	}

}

