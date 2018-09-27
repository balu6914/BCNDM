import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { Query } from '../../../common/interfaces/query.interface';

@Component({
  selector: 'dpc-sidebar-filters',
  templateUrl: './sidebar.filters.component.html',
  styleUrls: [ './sidebar.filters.component.scss' ]
})
export class SidebarFiltersComponent {
    _opened = false
    form: FormGroup;
    // TODO: Remove mocked cities
    mockCities: string[];
    @Output()
    // Emit event when we successfully buy more token , to get updated balance.
    filtersUpdate = new EventEmitter();
    public selectCityControl = new FormControl()

    constructor(
      private formBuilder: FormBuilder,
    ) {

       this.mockCities = ['Belgrade', 'Paris', 'Tel Aviv'];
       // Update form control city value on selected value change
       this.selectCityControl.valueChanges.subscribe(value =>
         this.form.patchValue({city: value})
       );

       this.form = formBuilder.group({
         name: ['', String],
         streamType: ['', String],
         minPrice: [0, Number],
         maxPrice: [1000000, Number]
       })
    }

    _toggleSidebar() {
      this._opened = !this._opened;
    }

    onClear() {
      this.form.reset();
    }

    onSubmit(isValid: boolean) {
    if(isValid) {
        this.filtersUpdate.emit(this.form.value)
      }
    }

}
