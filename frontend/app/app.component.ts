import { Component, OnInit } from '@angular/core';
import * as moment from 'moment/moment';
import { Developer, TimeConstraint } from './models'
import { CompanyService } from './company.service'

@Component({
  selector: 'my-app',
  template: `
    <h1>STATKUBE</h1>
	<time-picker [selected_time]="selected_time"></time-picker>
	<company-list [comps]="comps"></company-list>
	<dev-list [devs]="devs"></dev-list>
  `,
  providers: [CompanyService]
})

export class AppComponent implements OnInit{
  comps: Developer[];
  devs = DEVS;
  selected_time = new TimeConstraint("", "", "");

  constructor(private companyService: CompanyService) {};

  getCompanies(): void {
	  this.companyService.getCompanies().then(comps => this.comps = comps);
  }

  ngOnInit(): void {
	  this.getCompanies();
  }
}

const DEVS: Developer[] = [
    new Developer("super dev", 12),
    new Developer("super dev2", 13),
    new Developer("super dev3", 5),
    new Developer("super dev4", 120),
];


