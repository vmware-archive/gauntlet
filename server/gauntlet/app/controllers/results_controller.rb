class ResultsController < ApplicationController

  def index
    render status: 200, json: { results: Result.all }
  end

  def create

    result = Result.new(
      pipeline: params.fetch(:pipeline),
      pipecount: params.fetch('pipecount'),
      stage: params.fetch('stage'),
      stagecount: params.fetch('stagecount'),
      jobname: params.fetch('jobname'),
      gitinfo: params.fetch('gitinfo'),
      pass: params.fetch('pass')
    )

    if result.save
      render status: 201, json: { resultid: result.id }
    else
      render status: 400, json: { errors: resource.errors }
    end
  end


end
