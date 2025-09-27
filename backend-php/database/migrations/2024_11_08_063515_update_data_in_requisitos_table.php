<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class UpdateDataInRequisitosTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('requisitos', function (Blueprint $table) {
            $table->string('type_input')->default('text');
            $table->string('text_color')->default('default');
            $table->string('text_type')->default('normal');
            $table->string('text_size')->default('md');
            $table->boolean('is_recoverable')->default(false);
            $table->boolean('is_dependent')->default(false);
            $table->string('field_dependent')->default('');
            $table->string('value_dependent')->default('');
            $table->boolean('show_dependent')->default(false);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
{
    Schema::table('requisitos', function (Blueprint $table) {
        if (Schema::hasColumn('requisitos', 'type_input')) {
            $table->dropColumn('type_input');
        }
        if (Schema::hasColumn('requisitos', 'text_color')) {
            $table->dropColumn('text_color');
        }
        if (Schema::hasColumn('requisitos', 'text_type')) {
            $table->dropColumn('text_type');
        }
        if (Schema::hasColumn('requisitos', 'text_size')) {
            $table->dropColumn('text_size');
        }
        if (Schema::hasColumn('requisitos', 'is_recoverable')) {
            $table->dropColumn('is_recoverable');
        }
        if (Schema::hasColumn('requisitos', 'is_dependent')) {
            $table->dropColumn('is_dependent');
        }
        if (Schema::hasColumn('requisitos', 'field_dependent')) {
            $table->dropColumn('field_dependent');
        }
        if (Schema::hasColumn('requisitos', 'value_dependent')) {
            $table->dropColumn('value_dependent');
        }
        if (Schema::hasColumn('requisitos', 'show_dependent')) {
            $table->dropColumn('show_dependent');
        }
    });
}

}
