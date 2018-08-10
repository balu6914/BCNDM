import { ComponentRef, Injectable } from '@angular/core';
import { AlertContainerComponent } from '../alert-container/alert-container.component';
import { AlertComponent } from '../alert/alert.component';
import { AlertType } from '../alerts';
import { DomService } from './dom.service';

@Injectable()
export class AlertService {

  private alertDelay = 5000;
  private alertContainerElem: HTMLElement;
  private alertContainerRef: ComponentRef<AlertContainerComponent>;

  constructor(private domService: DomService) {
    this.alertContainerRef = this.domService.createRef(AlertContainerComponent);
    this.alertContainerElem = this.domService.getDomElementFromRef(this.alertContainerRef);
    this.domService.addChild(this.alertContainerElem);
  }

  success(message: string) {
    this.addAlert(message, AlertType.Success);
  }

  error(message: string) {
    this.addAlert(message, AlertType.Error);
  }

  warning(message: string) {
    this.addAlert(message, AlertType.Warning);
  }

  info(message: string) {
    this.addAlert(message, AlertType.Info);
  }

  private addAlert(message: string, type: AlertType) {
    const alertRef = this.domService.createRef(AlertComponent);
    alertRef.instance.message = message;
    alertRef.instance.type = type;
    alertRef.instance.ref = alertRef;
    const alertElement = this.domService.getDomElementFromRef(alertRef);
    this.domService.addChild(alertElement, this.alertContainerElem);
    this.domService.destroyAfter(alertRef, this.alertDelay);
  }

}
