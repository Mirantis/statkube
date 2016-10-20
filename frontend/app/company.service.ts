import { Injectable  } from '@angular/core';
import { Headers, Http  } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Developer } from './models'

@Injectable()
export class CompanyService {
	private companiesURL = 'http://172.18.163.4:8080/prstats/company'
	constructor(private http: Http) {};
	getCompanies(): Promise<Developer[]>{
		return this.http.get(this.companiesURL)
			.toPromise()
			.then(response => response.json().map(data => new Developer(data["FullName"], data["PRCount"])) as Developer[])
			.catch(this.handleError)
	};
	private handleError(error: any): Promise<any> {
		console.error('An error occurred', error); // for demo purposes only
		return Promise.reject(error.message || error);
	}
}
