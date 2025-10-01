import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostulationDynamicComponent } from './postulation-dynamic.component';

describe('PostulationDynamicComponent', () => {
  let component: PostulationDynamicComponent;
  let fixture: ComponentFixture<PostulationDynamicComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PostulationDynamicComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PostulationDynamicComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
