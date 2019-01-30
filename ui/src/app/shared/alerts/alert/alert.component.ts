import { animate, style, transition, trigger } from '@angular/animations';
import { Component, ComponentRef, Input } from '@angular/core';
import { AlertType } from 'app/shared/alerts/alerts';

@Component({
  selector: 'dpc-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss'],
  animations: [
    trigger(
      'alertAnimation', [
        transition('void => *', [
          style({
            transform: 'scale(0) translateX(15px)',
            opacity: 1
          }),
          animate('0.5s cubic-bezier(0.3, 1, 0.32, 1)')
        ]),
        transition('* => void', [
          animate('0.3s cubic-bezier(0.3, 1, 0.32, 1)', style({
            transform: 'scaleX(1) translateX(15px)',
            opacity: 0
          }))
        ])
      ]
    )
  ]
})
export class AlertComponent {
  @Input() message: string;
  @Input() type: AlertType;
  types = AlertType;
  ref: ComponentRef<any>;
  show = true;

  constructor() { }

  close() {
    this.show = false;
  }

  destroy() {
    if (!this.show) {
      this.ref.destroy();
    }
  }

}
