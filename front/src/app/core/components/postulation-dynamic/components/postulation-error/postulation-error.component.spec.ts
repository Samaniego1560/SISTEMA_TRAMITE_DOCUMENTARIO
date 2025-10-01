import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostulationErrorComponent } from './postulation-error.component';

describe('PostulationErrorComponent', () => {
  let component: PostulationErrorComponent;
  let fixture: ComponentFixture<PostulationErrorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PostulationErrorComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PostulationErrorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
