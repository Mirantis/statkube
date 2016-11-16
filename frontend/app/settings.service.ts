import { Injectable  } from '@angular/core';
import { Headers, Http  } from '@angular/http';

import 'rxjs/add/operator/toPromise';


@Injectable()
export class SettingsService{
	constructor(private http: Http) {};
	getSettings(): Promise<any>{
		return this.http.get("/settings.json").toPromise()
	};
}
