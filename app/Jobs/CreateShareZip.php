<?php

namespace App\Jobs;

use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use App\Models\Share;

class CreateShareZip implements ShouldQueue
{
    use Queueable;

    /**
     * Create a new job instance.
     */
    public function __construct(public Share $share)
    {
        //
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        try {
            //just check that we've not already created the zip file
            $zipPath = storage_path('app/shares/' . $this->share->user_id . '/' . $this->share->long_id . '.zip');
            if (file_exists($zipPath)) {
                return;
            }

            //if there is only one file just leave it alone and set the status to ready
            if ($this->share->file_count == 1) {
                $this->share->status = 'ready';
                $this->share->save();
                return;
            }

            //create the zip file
            $zip = new \ZipArchive();
            $zip->open($zipPath, \ZipArchive::CREATE | \ZipArchive::OVERWRITE);
            $files = $this->share->files;
            foreach ($files as $file) {
                $zip->addFile($this->share->path . '/' . $file->name, $file->name);
            }
            $zip->close();

            //if the zip file was created successfully, delete the files and folder and set the status to ready
            if (file_exists($zipPath)) {
                //delete the files and folder
                foreach ($files as $file) {
                    unlink($this->share->path . '/' . $file->name);
                }
                rmdir($this->share->path);

                //update the share status to ready
                $this->share->status = 'ready';
                $this->share->save();
            } else {
                //update the share status to failed
                $this->share->status = 'failed';
                $this->share->save();
            }
        } catch (\Exception $e) {
            $this->share->status = 'failed';
            $this->share->save();
            Log::error('Error creating share zip: ' . $e->getMessage());
        }
    }
}
