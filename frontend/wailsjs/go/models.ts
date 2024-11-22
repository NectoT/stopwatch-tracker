export namespace main {
	
	export enum ColorTheme {
	    Light = "light",
	    Dark = "dark",
	}
	export class StopwatchData {
	    IsActive: boolean;
	    Hue: number;
	    Name: string;
	    // Go type: time
	    LastUpdated: any;
	    // Go type: time
	    CreatedAt: any;
	    TimeAccumulated: number;
	
	    static createFrom(source: any = {}) {
	        return new StopwatchData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsActive = source["IsActive"];
	        this.Hue = source["Hue"];
	        this.Name = source["Name"];
	        this.LastUpdated = this.convertValues(source["LastUpdated"], null);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.TimeAccumulated = source["TimeAccumulated"];
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

