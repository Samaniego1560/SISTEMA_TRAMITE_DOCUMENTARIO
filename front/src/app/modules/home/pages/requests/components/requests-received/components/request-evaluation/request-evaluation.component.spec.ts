import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RequestEvaluationComponent } from './request-evaluation.component';

describe('RequestEvaluationComponent', () => {
  let component: RequestEvaluationComponent;
  let fixture: ComponentFixture<RequestEvaluationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RequestEvaluationComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RequestEvaluationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
