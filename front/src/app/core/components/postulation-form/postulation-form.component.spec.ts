import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostulationFormComponent } from './postulation-form.component';

describe('PostulationFormComponent', () => {
  let component: PostulationFormComponent;
  let fixture: ComponentFixture<PostulationFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PostulationFormComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PostulationFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
