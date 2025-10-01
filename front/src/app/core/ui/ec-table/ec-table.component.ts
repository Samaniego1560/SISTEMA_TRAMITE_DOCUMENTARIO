import {
  AfterContentInit,
  Component,
  ContentChildren,
  Input,
  QueryList,
  TemplateRef,
  ViewEncapsulation
} from '@angular/core';
import {EcColumnDef} from "./directives";
import {EcTemplate} from "../../directives/ec-template.directive";


@Component({
  selector: 'ec-table',
  templateUrl: './ec-table.component.html',
  styleUrls: ['./ec-table.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class EcTableComponent<T> implements AfterContentInit  {

  @Input() dataSource: T[] = []
  @Input() scrollHeight: string | undefined = undefined

  @ContentChildren(EcColumnDef) templatesColumn!: QueryList<EcColumnDef>;
  @ContentChildren(EcTemplate) templates!: QueryList<EcTemplate>;
  public captionTemplate?: TemplateRef<any>;
  public headerTemplate?: TemplateRef<any>;
  public bodyTemplate?: TemplateRef<any>;
  public emptyMessageTemplate?: TemplateRef<any>;

  constructor() {
  }

  ngAfterContentInit() {
    this.templates?.forEach((item) => {
      if (item.getType() === 'caption') {
        this.captionTemplate = item.template;
      }
      if (item.getType() === 'header') {
        this.headerTemplate = item.template;
      }
      if (item.getType() === 'body') {
        this.bodyTemplate = item.template;
      }
      if (item.getType() === 'emptymessage') {
        this.emptyMessageTemplate = item.template;
      }
    });
  }

  public getTemplateRow(header: EcColumnDef): TemplateRef<any> {
    if(header.cellDef === undefined){
      console.error('Not found column', header.getType())
    }
    return header.cellDef.template
  }
}
