import { NgFor, NgIf } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
    selector: 'app-modal-srq',
    standalone: true,
    imports: [NgFor, NgIf],
    templateUrl: './modal-srq.component.html',
    styleUrl: './modal-srq.component.scss'
})
export class ModalSrqComponent implements OnInit{

    @Input() participante: any; 
    @Input() itemSrq: any;
    @Input() questions: any;
    student: any;
    additionalInfoStudent: any
    dataAnswers: any = [];

    @Output() closeInfo = new EventEmitter<boolean>();

    ngOnInit(): void {
        this.student = this.itemSrq.data;
        this.additionalInfoStudent = this.participante;
        this.dataAnswers = this.questions;
    }

    exitInfoHistory(){
        this.closeInfo.emit(false);
    }

}
