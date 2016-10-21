import { Component, OnInit } from '@angular/core';

import * as moment from 'moment/moment';
import 'chartjs';

import { Developer, TimeConstraint } from './models'
import { CompanyService } from './company.service'

@Component({
  selector: 'my-app',
  template: `
    <h1>STATKUBE</h1>
	<time-picker [selected_time]="selected_time"></time-picker>
	<button type="button" class="btn btn-primary" (click)="filter()">Filter by date</button>
	<company-list [comps]="comps"></company-list>
	<!--dev-list [devs]="devs"></dev-list-->
  `,
  providers: [CompanyService]
})

export class AppComponent implements OnInit{
  comps: Developer[];
  devs = DEVS;
  selected_time = new TimeConstraint("", "", "");

  constructor(private companyService: CompanyService) {};

  getCompanies(): void {
	  this.companyService.getCompanies(this.selected_time).then(comps => this.comps = comps);
  }

  filter(): void{
	  this.getCompanies();
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


