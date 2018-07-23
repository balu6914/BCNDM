import { Component } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import { Router } from '@angular/router';

import { StreamService } from '../services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';

@Component({
  selector: 'dashboard-sell-add',
  templateUrl: './dashboard.sell.add.component.html',
  styleUrls: [ './dashboard.sell.add.component.scss' ]
})
export class DashboardSellAddComponent {
    user:any;

    public form: FormGroup;
    constructor(
        private router: Router,
        private fb: FormBuilder,
        private StreamService: StreamService,
    ) {
    }
    ngOnInit() {
        this.form = this.fb.group({
                name       : ['', [<any>Validators.required, <any>Validators.minLength(3)]],
                type       : ['', [<any>Validators.required]],
                description: ['', [<any>Validators.required, <any>Validators.minLength(5)]],
                url        : ['', [<any>Validators.required]],
                price      : ['', [<any>Validators.required]],
                long       : ['', [<any>Validators.required]],
                lat        : ['', [<any>Validators.required]]
        });
    }

    onSubmit(model: Stream, isValid: boolean) {
        if(isValid) {
            // Confirm dialog
            let confirmMsg = 'Your Stream will be published automatically on the market, affter creation. Do you Agree ?'
            // let result = this.dialogService.confirm(confirmMsg, 'Cancel', 'Yes, publish it!');
            //  result.subscribe( () => {
            //      console.log('confirmed');
            //      // Convert to miTAS
            //      var mitas = model["price"] * Math.pow(10, 6);
            //      model["price"] = mitas;
            //      model["location"] = {
            //          "type": "Point",
            //          "coordinates": [parseFloat(model["long"]),
            //                          parseFloat(model["lat"])]
            //                      }
            //     delete model["long"];
            //     delete model["lat"];
            //
            //      this.StreamService.addStream(model).subscribe(
            //          response => {
            //              this.router.navigate(['/dashboard/sell/map'])
            //          },
            //          err => {
            //              console.log(err);
            //          });
            //    },
            //    (err: any) => {
            //      console.log('declined');
            //    }
            //  );
            //  // if you only need the confirm answer
            //  result.onErrorResumeNext().subscribe( () => {
            //    console.log('confirmed 2');
            //  });
        }
    }
}
