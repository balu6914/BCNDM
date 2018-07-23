import { Component } from '@angular/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';

import { StreamService } from '../services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { TasPipe, MitasPipe } from '../../../common/pipes/converter.pipe';

@Component({
  selector: 'dashboard-sell-edit',
  templateUrl: './dashboard.sell.edit.component.html',
  styleUrls: [ './dashboard.sell.edit.component.scss' ]
})
export class DashboardSellEditComponent {
    user:any;
    stream: any;
    editName: string;
    editType: string;
    editDescription: string;
    editUrl: string;
    public form: FormGroup;

    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private fb: FormBuilder,
        private StreamService: StreamService,
        private tasPipe: TasPipe,
        private mitasPipe: MitasPipe,
    ) {
    }
    ngOnInit() {
        // Fetch stream
        let id = this.route.snapshot.params['id'];
        this.StreamService.getStream(id).subscribe(
                (result: any) => {
                    this.stream = result.Stream;
                    // Convert to TAS
                    this.stream["price"] = this.tasPipe.transform(this.stream["price"]);
                    // Pre-set FormGroup with stream values
                    this.form.controls['name'].setValue(this.stream["name"]);
                    this.form.controls['type'].setValue(this.stream["type"]);
                    this.form.controls['description'].setValue(this.stream["description"]);
                    this.form.controls['url'].setValue(this.stream["url"]);
                    this.form.controls['price'].setValue(this.stream["price"]);
                    this.form.controls['long'].setValue(this.stream["location"]["coordinates"][0]);
                    this.form.controls['lat'].setValue(this.stream["location"]["coordinates"][1]);
                },
                err => { console.log(err) }
              );

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
            let confirmMsg = 'Your Stream will be published automatically on the market affter edition. Do you Agree ?'
            // let result = this.dialogService.confirm(confirmMsg, 'Cancel', 'Yes, publish it!');
            // result.subscribe( () => {
            //     // Convert to miTAS
            //     model["price"] = this.mitasPipe.transform(model["price"]);
            //     model["location"] = {
            //         "type": "Point",
            //         "coordinates": [parseFloat(model["long"]),
            //                         parseFloat(model["lat"])]
            //     }
            //     delete model["long"];
            //     delete model["lat"];
            //
            //     this.StreamService.updateStream(this.stream["id"], model).subscribe(
            //         response => {
            //             console.log(response);
            //         },
            //         err => {
            //             if (err.status == 200) {
            //                 this.router.navigate(['/dashboard/sell/map'])
            //             } else {
            //                 console.log(err);
            //             }
            //         });
            //     },
            //     (err: any) => {
            //         console.log('declined');
            //     }
            // );
            // // if you only need the confirm answer
            // result.onErrorResumeNext().subscribe( () => {
            //     console.log('confirmed 2');
            // });
        }
    }

    onDelete() {
        let id = this.route.snapshot.params['id'];
        // Confirm dialog
        let confirmMsg = 'Your Stream will be removed from the market. Do you Agree ?'
        // let result = this.dialogService.confirm(confirmMsg, 'Cancel', 'Yes, remove it!');
        //  result.subscribe( () => {
        //      console.log('confirmed');
        //      this.StreamService.removeStream(id).subscribe(
        //          response => {
        //              console.log(response);
        //          },
        //          err => {
        //              if (err.status == 200) {
        //                  this.router.navigate(['/dashboard/sell/map'])
        //              } else {
        //                  console.log(err);
        //              }
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
