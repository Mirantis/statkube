import { Injectable  } from '@angular/core';
import { Developer } from './models'

@Injectable()
export class CompanyService {
	getCompanies(): Promise<Developer[]>{
		return Promise.resolve(COMPANIES);
	}
}

const COMPANIES: Developer[] = [
    new Developer("Manimuru", 45561),
    new Developer("CoreOS", 2),
    new Developer("GlueGL", 444),
]
