import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormArray, FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';

import { AuthService } from 'app/auth/services/auth.service';
import { User } from 'app/common/interfaces/user.interface';
import { Query } from 'app/common/interfaces/query.interface';
import { AccessService } from 'app/dashboard/access/access.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-access-add',
  templateUrl: './dashboard.access.add.component.html',
  styleUrls: [ './dashboard.access.add.component.scss' ]
})
export class DashboardAccessAddComponent implements OnInit {
  user = <User>{};
  query = new Query();
  form: FormGroup;
  parties: FormArray;
  usersName = [];
  streams = [];
  submitted = false;

  users = [];

  @Output() accessCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private authService: AuthService,
    private accessService: AccessService,
    private formBuilder: FormBuilder,
    public modalNewAccess: BsModalRef,
    public alertService: AlertService,
  ) {

    this.form = this.formBuilder.group({
      parties:    this.formBuilder.array([this.createPartner()])
    }, {
      validator: [this.partnerValidator.bind(this)]
    });

  }

  createPartner(): FormGroup {
    return this.formBuilder.group({
      partner: ['', [Validators.required]],
    });
  }

  partnerValidator(fg: FormGroup) {
    // Verify that partner ID is not the owner ID
    fg.value.parties.forEach( (item, i) => {
        if (item.partner === this.user.id) {
          // Create a custom error field used as *ngIf condition for style
          fg.controls.parties['controls'][i].controls.partner.setErrors({
            'ownerID': true
          });
        }
    });
  }

  ngOnInit () {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;

        // TODO: Fetch users
        this.users = [
          { name: 'Aeroflot',                 id: '1'},
          { name: 'Aerolíneas Argentinas',    id: '2'},
          { name: 'Aeroméxico',               id: '3'},
          { name: 'Aeroflot',                 id: '4'},
          { name: 'Air Europa',               id: '5'},
          { name: 'Air France',               id: '6'},
          { name: 'Alitalia',                 id: '7'},
          { name: 'China Airlines',           id: '8'},
          { name: 'China Eastern Airlines',   id: '9'},
          { name: 'CSA Czech Airlines',       id: '10'},
          { name: 'Delta Air Lines',          id: '11'},
          { name: 'Garuda Indonesia',         id: '12'},
          { name: 'Kenya Airways',            id: '13'},
          { name: 'KLM Royal Dutch Airlines', id: '14'},
          { name: 'Korean Air',               id: '15'},
          { name: 'Middle East Airlines',     id: '16'},
          { name: 'Saudia',                   id: '17'},
          { name: 'Tarom',                    id: '18'},
          { name: 'Vietnam Airlines',         id: '19'},
          { name: 'Xiamen Air',               id: '20'},
        ];

        // Fetch Datapace users
        this.usersName = this.users.map(stream => stream.name);
      },
      err => {
        console.log(err);
      }
    );
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {

      const createAccessReq = {
        items: [],
      };

      this.form.value.parties.forEach( item => {
        const partner = {
          partner_id: item.partner,
        };
        createAccessReq.items.push(partner);
      });

      // TODO: Send access request
      this.accessCreated.emit(createAccessReq);

      this.modalNewAccess.hide();
    }
  }

  onAddPartner() {
    this.parties = this.form.get('parties') as FormArray;
    this.parties.push(this.createPartner());
  }

  onDeletePartner(index: number) {
    this.parties.controls.splice(index, 1);
  }
}
