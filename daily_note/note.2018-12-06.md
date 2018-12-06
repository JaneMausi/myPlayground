### when command yarn install does not work in Git bash, what to do?

I got this same error. I was able to reproduce and fix the issue (it was direct mistake from me).

Basically, I have a git dependency (and it's pointing a branch) that is tracked in the yarn lockfile. At some point during the day, I force pushed my git branch which clobbered the tracked sha. I verified that the clobbered sha was the git archive CLOBBERED_SHA presented in the error message.

The fix was to yarn upgrade <git-dependency> (if you're on latest). Otherwise, you can yarn add <git+ssh://...#branch> --force.

### rename all files in same directory by adding suffix

	$ eval $(echo $(ls -rlt |awk 'V==0{V=1;next;}{printf("mv %s %s.md;",$NF,$NF)}'))