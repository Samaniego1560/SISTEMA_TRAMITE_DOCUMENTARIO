<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateThesisPracticacingsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('thesis_practicacings', function (Blueprint $table) {
            $table->id();
            $table->string('type_student');
            $table->string('cod_student');
            $table->string('name');
            $table->string('appaterno');
            $table->string('apmaterno');
            $table->string('sexo');
            $table->string('age');
            $table->string('facultad');
            $table->string('escuela_profesional');
            $table->string('mod_ingreso');
            $table->string('dni');
            $table->string('email');
            $table->unsignedBigInteger('convocatoria_id');
            $table->unsignedBigInteger('status_id');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('thesis_practicacings');
    }
}
