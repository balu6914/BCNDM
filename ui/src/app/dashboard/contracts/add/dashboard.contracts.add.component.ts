import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import { MdlDatePickerService } from '@angular-mdl/datepicker';
import { MdlPopoverComponent } from '@angular2-mdl-ext/popover';
import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'dashboard-contracts-add',
  templateUrl: './dashboard.contracts.add.component.html',
  styleUrls: [ './dashboard.contracts.add.component.scss' ]
})
export class DashboardContractsAddComponent {

    public form: FormGroup;
    // TODO: Remove this Mock contracts list for demo purpose
    streamList = [
        {'id' : 1, 'stream': {'name': 'WeIO pressure', 'price':'1'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':false, 'expiered': false},
        {'id' : 2, 'stream': {'name': 'WeIO temperature', 'price':'10'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':true, 'expiered': false },
        {'id' : 3, 'stream': {'name': 'WeIO humidity', 'price':'15'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': false},
        {'id' : 4, 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':'true', 'expiered': false},
        {'id' : 5, 'stream': {'name': 'WeIO radiation', 'price':'50'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': false},
        {'id' : 5, 'stream': {'name': 'Spark', 'price':'30'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': true},
    ];
    usersList = [
        {'id' : 1, 'name': 'John Down', 'email':"john@castleblack.com"},
        {'id' : 1, 'name': 'Nikola Marcetic', 'email':"daenerys@fire.com"},
        {'id' : 1, 'name': 'Jason Born', 'email':"cersei@castlerock.com"},
    ];

    selectedEndDate: Date;
    parties: any;

    constructor(
        private AuthService: AuthService,
        private router: Router,
        private fb: FormBuilder,
        private datePicker: MdlDatePickerService
    ) { }

    ngOnInit() {
        // Hack to fix a MdlPopover firing TypeError on ngOnDestroy() angular2-mdl-ext issue
        MdlPopoverComponent.prototype.ngOnDestroy = function () {
            this.elementRef.nativeElement.removeEventListener(this, 'hide');
        };
        this.form = this.fb.group({
            stream       : ['', [<any>Validators.required]],
            exparationDate  : ['', [<any>Validators.required]],
            parties: ['', [<any>Validators.required]]
        });
    }

    // Datepicker handler
    public pickADate($event: MouseEvent) {
    this.datePicker.selectDate(this.selectedEndDate, {openFrom: $event}).subscribe( (selectedDate: Date) => {
        this.selectedEndDate = selectedDate;
    });
 }


}
