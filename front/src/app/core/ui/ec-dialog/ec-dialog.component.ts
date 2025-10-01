import {
  AfterContentInit,
  Component,
  ContentChildren,
  input,
  output,
  QueryList,
  TemplateRef,
} from '@angular/core';
import { NgIf, NgTemplateOutlet } from '@angular/common';
import {EcTemplate} from "../../directives/ec-template.directive";

@Component({
  selector: 'ec-dialog',
  standalone: true,
  imports: [NgIf, NgTemplateOutlet],
  templateUrl: './ec-dialog.component.html',
  styleUrl: './ec-dialog.component.scss',
})
export class EcDialogComponent implements AfterContentInit {
  header = input<string>();
  visible = input<boolean>();
  visibleChange = output<boolean>();

  style = input<any>();
  @ContentChildren(EcTemplate) templates!: QueryList<EcTemplate>;
  headlessTemplate?: TemplateRef<any>;
  headerTemplate?: TemplateRef<any>;
  footerTemplate?: TemplateRef<any>;

  ngAfterContentInit() {
    this.templates?.forEach((item) => {
      if (item.getType() === 'headless') {
        this.headlessTemplate = item.template;
      }
      if (item.getType() === 'header') {
        this.headerTemplate = item.template;
      }
      if (item.getType() === 'footer') {
        this.footerTemplate = item.template;
      }
    });
  }
}
