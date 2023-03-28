import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { Query } from 'app/common/interfaces/query.interface';

@Component({
  selector: 'dpc-sidebar-filters',
  templateUrl: './sidebar.filters.component.html',
  styleUrls: [ './sidebar.filters.component.scss' ]
})
export class SidebarFiltersComponent {
    _opened = false;
    form: FormGroup;
    submitted = false;

    @Output()
    // Emit event when we successfully buy more token , to get updated balance.
    filtersUpdate = new EventEmitter();
    public selectCityControl = new FormControl();

    constructor(
      private formBuilder: FormBuilder,
    ) {

       // Update form control city value on selected value change
       this.selectCityControl.valueChanges.subscribe(value =>
         this.form.patchValue({city: value})
       );

       this.form = formBuilder.group({
         name:  ['', [ Validators.maxLength(32)]],
         streamType: ['', [ Validators.maxLength(32)]],
         minPrice: ['', [ Validators.min(0.000001)]],
         maxPrice: ['', [ Validators.max(1000000)]]
       });
    }

    _toggleSidebar() {
      this._opened = !this._opened;
    }

    onClear() {
      this.form.reset();
      this.onSubmit();
    }

    onSubmit() {
      this.submitted = true;
      if (this.form.valid) {
        this.filtersUpdate.emit(this.form.value);
      }
    }
}
