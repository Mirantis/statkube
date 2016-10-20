import { Component, Input } from '@angular/core';
import * as moment from 'moment/moment';
import { Developer } from './models'

@Component({
  selector: 'dev-list',
  template: `
  <div>
    <h2>Engineer statistics</h2>
    <table class="table table-striped table-bordered">
      <thead>
        <tr>
          <th>Name</th>
          <th>PR Count</th>
        </tr>
        </thead>
      <tbody>
        <tr *ngFor="let dev of devs">
          <td>{{dev.name}}</td>
          <td>{{dev.pr_count}}</td>
        </tr>
      </tbody>
    </table>
  </div>
  `
})

export class DevListComponent{
  @Input()
  devs: Developer[];
}

