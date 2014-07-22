class CreateResults < ActiveRecord::Migration
  def change
    create_table :results do |t|
      t.string :pipeline
      t.integer :pipecount
      t.string :stage
      t.integer :stagecount
      t.string :jobname
      t.string :gitinfo
      t.boolean :pass

      t.timestamps
    end
  end
end
