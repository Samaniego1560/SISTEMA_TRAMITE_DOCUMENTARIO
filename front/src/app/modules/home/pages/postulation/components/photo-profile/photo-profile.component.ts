import {Component, EventEmitter, input, Output} from '@angular/core';
import {NgIf} from "@angular/common";
import {IRequirement} from "../../../../../../core/models/announcement";

@Component({
  selector: 'app-photo-profile',
  standalone: true,
    imports: [
        NgIf
    ],
  templateUrl: './photo-profile.component.html',
  styleUrl: './photo-profile.component.scss'
})
export class PhotoProfileComponent {
  requirement = input.required<IRequirement>();
  src = input.required<string>();
  name = input.required<string>();
  description = input.required<string>();
  @Output() onProcessFile = new EventEmitter<{event: any, key: string, type: number}>();
  public isUploading: boolean = false;

  onFileSelected($event: Event) {
    const key = this.requirement().id?.toString();
    this.onProcessFile.emit({event: $event, key: key || '', type: this.requirement().tipo_requisito_id})
    this.loadFile($event);
  }

  loadFile(event: any) {
    const file = event.target.files[0];
    if (file && file.type.startsWith('image/')) {
      this.isUploading = true;
      const reader = new FileReader();

      reader.onload = (e: any) => {
        //this.requirement().default = e.target.result;
        this.isUploading = false;
      };

      reader.readAsDataURL(file);
    }
  }
}
