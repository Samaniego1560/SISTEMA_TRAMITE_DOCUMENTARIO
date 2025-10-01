import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PreQuestionsComponent } from './pre-questions.component';

describe('PreQuestionsComponent', () => {
  let component: PreQuestionsComponent;
  let fixture: ComponentFixture<PreQuestionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PreQuestionsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PreQuestionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
