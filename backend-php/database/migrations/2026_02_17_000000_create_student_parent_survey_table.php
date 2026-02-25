<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateStudentParentSurveyTable extends Migration
{
    public function up()
    {
        Schema::create('student_parent_survey', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->string('student_code', 20);
            $table->string('dni', 12)->unique();
            $table->string('school', 150);
            $table->string('academic_cycle', 20);
            $table->unsignedInteger('age');
            $table->enum('gender', ['F', 'M']);

            $table->boolean('is_pregnant');
            $table->boolean('has_children');

            $table->unsignedTinyInteger('children_count')->nullable()->default(0);
            $table->unsignedTinyInteger('child_1_age')->nullable();
            $table->unsignedTinyInteger('child_2_age')->nullable();
            $table->unsignedTinyInteger('child_3_age')->nullable();
            $table->unsignedTinyInteger('child_4_age')->nullable();

            $table->enum('residence_type', ['familiar', 'alquiler', 'residencia', 'otro'])->nullable();
            $table->string('residence_other', 255)->nullable();

            $table->boolean('lives_with_partner')->nullable();
            $table->boolean('lives_with_children')->nullable();

            $table->enum('caregiver_type', ['ambos', 'familia', 'padre', 'madre', 'otro'])->nullable();
            $table->string('caregiver_other', 255)->nullable();

            $table->boolean('participates_in_upbringing')->nullable();
            $table->boolean('economic_support')->nullable();
            $table->decimal('monthly_amount', 10, 2)->nullable();

            $table->boolean('has_academic_support')->nullable();
            $table->enum('support_provider', ['ambos', 'familia', 'padre', 'madre', 'otro'])->nullable();
            $table->string('support_provider_other', 255)->nullable();

            $table->boolean('interrupted_semester')->nullable();

            $table->timestamps();
        });
    }

    public function down()
    {
        Schema::dropIfExists('student_parent_survey');
    }
}
