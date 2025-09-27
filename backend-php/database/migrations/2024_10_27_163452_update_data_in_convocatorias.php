<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class UpdateDataInConvocatorias extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('convocatorias', function (Blueprint $table) {
            $table->integer('credito_minimo')->nullable();
            $table->integer('nota_minima')->nullable();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('convocatorias', function (Blueprint $table) {
            $table->dropColumn('credito_minimo');
            $table->dropColumn('nota_minima');
        });
    }
}
