import { Injectable  } from '@angular/core';
import { Headers, Http  } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Developer, TimeConstraint } from './models'
import { SettingsService } from './settings.service'

@Injectable()
export class CompanyService {
	private companiesPath = '/prstats/company'
	constructor(private http: Http) {};
	getCompanies(time: TimeConstraint, settings: any): Promise<Developer[]>{
		return settings.then(response => {
			return response.json()["api_root"]
		})
		.then(root => {
			var url = root + this.companiesPath + "?" + time.getParams();
			return this.http.get(url)
				.toPromise()
				.then(response => response.json()
						.filter(data => data["PRCount"] > 0)
						.map(data => new Developer(data["FullName"], data["PRCount"])) as Developer[])
				.catch(this.handleError)
		})
		.catch(this.handleError)
	};

	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error); // for demo purposes only
		return Promise.reject(error.message || error);
	}
}
