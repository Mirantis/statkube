import { Injectable  } from '@angular/core';
import { Headers, Http  } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Developer, TimeConstraint } from './models'

@Injectable()
export class CompanyService {
	private companiesURL = 'http://172.18.163.4:8080/prstats/company'
	constructor(private http: Http) {};
	getCompanies(time: TimeConstraint): Promise<Developer[]>{
		var url = this.companiesURL + "?" + time.getParams();
		return this.http.get(url)
			.toPromise()
			.then(response => response.json()
                    .filter(data => data["PRCount"] > 0)
                    .map(data => new Developer(data["FullName"], data["PRCount"])) as Developer[])
			.catch(this.handleError)
	};

	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error); // for demo purposes only
		return Promise.reject(error.message || error);
	}
}
