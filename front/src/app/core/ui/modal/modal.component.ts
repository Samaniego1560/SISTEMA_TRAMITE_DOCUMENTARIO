import {Component, EventEmitter, Input, Output, ElementRef, ViewChild, AfterViewInit, OnChanges, SimpleChanges} from '@angular/core';
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-modal',
  standalone: true,
  imports: [
    NgIf
  ],
  templateUrl: './modal.component.html',
  styleUrl: './modal.component.scss'
})
export class ModalComponent implements OnChanges, AfterViewInit {

  @Input() public show: boolean = false;
  @Output() public showChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  @Output() onNext: EventEmitter<void> = new EventEmitter<void>();
  @Input() public title: string = '';

  @Input() showFooter: boolean = true;
  @Input() showHeader: boolean = true;
  @Input() showClose: boolean = true;

  @Input() buttonSuccess: string = 'Guardar';
  @Input() buttonCancel: string = 'Cancelar';

  @Input() showButtonCancel: boolean = false;

  @ViewChild('modalContainer', { static: false }) modalContainer?: ElementRef;

  private previouslyFocusedElement?: HTMLElement;

  ngAfterViewInit() {
    // Handle initial focus if modal is already open
    if (this.show) {
      this.handleModalOpen();
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['show']) {
      if (changes['show'].currentValue) {
        this.handleModalOpen();
      } else {
        this.handleModalClose();
      }
    }
  }

  private handleModalOpen() {
    // Store the currently focused element to restore later
    this.previouslyFocusedElement = document.activeElement as HTMLElement;
    
    // Focus the modal after a short delay to ensure it's rendered
    setTimeout(() => {
      if (this.modalContainer) {
        this.modalContainer.nativeElement.focus();
      }
    }, 100);
  }

  private handleModalClose() {
    // Restore focus to the previously focused element
    if (this.previouslyFocusedElement && this.previouslyFocusedElement.focus) {
      this.previouslyFocusedElement.focus();
    }
  }

}
