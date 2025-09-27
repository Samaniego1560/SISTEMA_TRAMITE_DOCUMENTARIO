<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class UpdateDataInTipoRequisitosTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {

        DB::table('tipo_requisitos')->insert([
            [
                'id' => 5,
                'descripcion_tipo' => 'Radio'
            ],
            [
                'id' => 6,
                'descripcion_tipo' => 'Checkbox'
            ]
        ]);
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('tipo_requisitos', function (Blueprint $table) {
        });
    }
}
