import { Injectable, Injector, ComponentFactoryResolver, EmbeddedViewRef, ApplicationRef, ComponentRef } from '@angular/core';

@Injectable()
export class DomService {

  constructor(
    private componentFactoryResolver: ComponentFactoryResolver,
    private appRef: ApplicationRef,
    private injector: Injector
  ) { }

  createRef(component: any): ComponentRef<any> {
    const componentRef = this.componentFactoryResolver
      .resolveComponentFactory(component)
      .create(this.injector);
    this.appRef.attachView(componentRef.hostView);
    return componentRef;
  }

  getDomElementFromRef(componentRef: ComponentRef<any>): HTMLElement {
    return (componentRef.hostView as EmbeddedViewRef<any>)
      .rootNodes[0] as HTMLElement;
  }

  addChild(child: HTMLElement, parent: HTMLElement = document.body) {
    parent.appendChild(child);
  }

  destroyAfter(componentRef: ComponentRef<any>, delay: number) {
    setTimeout(() => {
        componentRef.instance.show = false;
    }, delay);
  }

}
