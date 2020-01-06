import { Component, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { environment } from 'environments/environment';

@Component({
  selector: 'dpc-kubeflow',
  templateUrl: './kubeflow.component.html',
  styleUrls: ['./kubeflow.component.scss'],
})
export class KubeflowComponent implements OnInit {
  iframeGrafana: any;

  constructor(
    private domSanitizer: DomSanitizer,
  ) { }

  ngOnInit() {
    this.iframeGrafana = this.domSanitizer.bypassSecurityTrustResourceUrl(environment.KUBEFLOW_URL);
  }
}
