import { Component, Input } from '@angular/core';
import * as moment from 'moment/moment';
import { Developer } from './models'

@Component({
  selector: 'company-list',
  template: `
  <div>
    <h2>Company statistics</h2>
    <table class="table table-striped table-bordered">
        <thead>
        <tr>
            <th>Name</th>
            <th>PR Count</th>
        </tr>
        </thead>
        <tbody>
        <tr *ngFor="let i of comps">
            <td>{{i.name}}</td>
            <td>{{i.pr_count}}</td>
        </tr>
        </tbody>
    </table>
  </div>
  `
})

export class CompanyListComponent{
  @Input()
  comps: Developer[];
}

