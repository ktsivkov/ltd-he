export namespace history {
	
	export class AppendRequest {
	    elo: number;
	    mvp: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppendRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.elo = source["elo"];
	        this.mvp = source["mvp"];
	    }
	}
	export class GameHistory {
	    outcome: string;
	    eloDiff: number;
	    date: string;
	    isLast: boolean;
	    account?: player.Player;
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
	    gameVersion: string;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new GameHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outcome = source["outcome"];
	        this.eloDiff = source["eloDiff"];
	        this.date = source["date"];
	        this.isLast = source["isLast"];
	        this.account = this.convertValues(source["account"], player.Player);
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
	        this.gameVersion = source["gameVersion"];
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

}

export namespace player {
	
	export class Player {
	    battleTag: string;
	    logsPathAbsolute: string;
	    logsPathRelative: string;
	    reportFilePathAbsolute: string;
	    reportFilePathRelative: string;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battleTag = source["battleTag"];
	        this.logsPathAbsolute = source["logsPathAbsolute"];
	        this.logsPathRelative = source["logsPathRelative"];
	        this.reportFilePathAbsolute = source["reportFilePathAbsolute"];
	        this.reportFilePathRelative = source["reportFilePathRelative"];
	    }
	}

}

