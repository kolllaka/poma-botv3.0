import { durationFormat } from '/static/js/lib/misc.js';

export function getTemplateSong(song, index) {
	return `
	<li class="playlist__item itemplaylist">
		<div class="itemplaylist__body">
			<div class="itemplaylist__cell">${index + 1}</div>
			<div class="itemplaylist__cell itemplaylist__cell-n">${song.name}</div>
			<div class="itemplaylist__cell">${durationFormat(song.duration)}</div>
			<div class="itemplaylist__cell">
				<div data-index="${index}" data-btn="del" class="btn"></div>
			</div>
		</div>
	</li>
	`
}

export function getTemplateInfo(info) {
	return `
	<div class="playlist__info">
		общее количество треков: ${info.count} на ${durationFormat(info.duration)}
	</div>
	`
}